package main

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
    "sync"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/PuerkitoBio/goquery"
    "github.com/redis/go-redis/v9"
)

type Crawler struct {
    redisClient *redis.Client
    s3Client    *s3.Client
    dynamoClient *dynamodb.Client
    httpClient  *http.Client
    rateLimiter map[string]*time.Ticker
    mu          sync.Mutex
}

type PageData struct {
    URL      string    `dynamodbav:"url"`
    Title    string    `dynamodbav:"title"`
    Links    []string  `dynamodbav:"links"`
    FetchedAt time.Time `dynamodbav:"fetched_at"`
}

func NewCrawler() (*Crawler, error) {
    // Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // AWS config
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        return nil, err
    }

    // HTTP client with timeout
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    return &Crawler{
        redisClient:  rdb,
        s3Client:     s3.NewFromConfig(cfg),
        dynamoClient: dynamodb.NewFromConfig(cfg),
        httpClient:   httpClient,
        rateLimiter:  make(map[string]*time.Ticker),
    }, nil
}

func (c *Crawler) EnqueueURL(ctx context.Context, urlStr string) error {
    // Normalize URL
    u, err := url.Parse(urlStr)
    if err != nil {
        return err
    }
    urlStr = u.String()

    // Check if already seen (cache hit)
    if c.redisClient.SIsMember(ctx, "seen_urls", urlStr).Val() {
        return nil // Skip duplicates
    }

    // Add to Redis queue
    return c.redisClient.LPush(ctx, "url_queue", urlStr).Err()
}

func (c *Crawler) Crawl(ctx context.Context) {
    for {
        // Pop URL from queue
        urlStr, err := c.redisClient.RPop(ctx, "url_queue").Result()
        if err == redis.Nil {
            time.Sleep(1 * time.Second) // Queue empty, wait
            continue
        }
        if err != nil {
            log.Printf("Error popping URL: %v", err)
            continue
        }

        // Mark as seen
        c.redisClient.SAdd(ctx, "seen_urls", urlStr)

        // Fetch and process
        if err := c.processURL(ctx, urlStr); err != nil {
            log.Printf("Error processing %s: %v", urlStr, err)
        }
    }
}

func (c *Crawler) processURL(ctx context.Context, urlStr string) error {
    // Rate limiting per domain
    u, _ := url.Parse(urlStr)
    domain := u.Hostname()
    c.mu.Lock()
    if _, ok := c.rateLimiter[domain]; !ok {
        c.rateLimiter[domain] = time.NewTicker(1 * time.Second) // 1 req/sec per domain
    }
    ticker := c.rateLimiter[domain]
    c.mu.Unlock()
    <-ticker.C

    // Fetch page
    resp, err := c.httpClient.Get(urlStr)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("status %d", resp.StatusCode)
    }

    // Cache raw content in Redis (short TTL)
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    c.redisClient.Set(ctx, "cache:"+urlStr, body, 1*time.Hour)

    // Parse HTML
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
    if err != nil {
        return err
    }

    // Extract data
    data := c.parsePage(doc, urlStr)

    // Store raw page in S3
    hash := md5.Sum([]byte(urlStr))
    key := hex.EncodeToString(hash[:])
    _, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String("web-crawler-data"),
        Key:    aws.String(key),
        Body:   bytes.NewReader(body),
    })
    if err != nil {
        return err
    }

    // Store metadata in DynamoDB
    av, err := attributevalue.MarshalMap(data)
    if err != nil {
        return err
    }
    _, err = c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: aws.String("WebCrawlerMetadata"),
        Item:      av,
    })
    if err != nil {
        return err
    }

    // Enqueue new URLs
    for _, link := range data.Links {
        c.EnqueueURL(ctx, link)
    }

    return nil
}

func (c *Crawler) parsePage(n *goquery.Document, baseURL string) PageData {
    var data PageData
    data.URL = baseURL
    data.FetchedAt = time.Now()

    var links []string
    n.Find("a").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        if exists {
            absURL, err := url.Parse(baseURL).Parse(href)
            if err == nil && (absURL.Scheme == "http" || absURL.Scheme == "https") {
                links = append(links, absURL.String())
            }
        }
    })
    data.Links = links

    // Extract title
    data.Title = n.Find("title").Text()

    return data
}

func main() {
    crawler, err := NewCrawler()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    crawler.EnqueueURL(ctx, "https://example.com")
    go crawler.Crawl(ctx)

    select {} // Run indefinitely
}