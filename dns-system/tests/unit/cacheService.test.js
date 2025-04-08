const { getFromCache, setInCache } = require('../../backend/services/cacheService');
const redis = require('redis');

jest.mock('redis', () => ({
  createClient: jest.fn().mockReturnValue({
    get: jest.fn(),
    setex: jest.fn(),
  }),
}));

describe('cacheService', () => {
  it('should get from cache', async () => {
    const client = require('redis').createClient();
    client.get.mockImplementation((key, cb) => cb(null, JSON.stringify('192.168.1.1')));
    const result = await getFromCache('example.com');
    expect(result).toBe('192.168.1.1');
  });
});