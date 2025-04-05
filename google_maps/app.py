from flask import Flask, render_template, request, jsonify
import folium
import networkx as nx
from math import radians, cos, sin, asin, sqrt
import queue
import threading
import time

app = Flask(__name__)

# --- Location Search Service ---
class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end_of_word = False

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word):
        node = self.root
        for char in word.lower():
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end_of_word = True

    def search_prefix(self, prefix):
        node = self.root
        for char in prefix.lower():
            if char not in node.children:
                return []
            node = node.children[char]
        return self._find_words_from_node(node, prefix.lower())

    def _find_words_from_node(self, node, current_prefix):
        words = []
        if node.is_end_of_word:
            words.append(current_prefix)
        for char, child in node.children.items():
            words.extend(self._find_words_from_node(child, current_prefix + char))
        return words

locations = {
    'Los Angeles': (34.0522, -118.2437),
    'Las Vegas': (36.1699, -115.1398),
    'San Francisco': (37.7749, -122.4194),
    'New York': (40.7128, -74.0060),
    'Los Alamos': (35.8800, -106.3031),
    'Las Cruces': (32.3199, -106.7637),
    'Chicago': (41.8781, -87.6298),
    'Miami': (25.7617, -80.1918),
    'Seattle': (47.6062, -122.3321),
    'Denver': (39.7392, -104.9903)
}

location_trie = Trie()
for loc in locations:
    location_trie.insert(loc)

def search_location(query):
    suggestions = location_trie.search_prefix(query)
    if not suggestions:
        return None, "No locations found"
    exact_match = next((s for s in suggestions if s.lower() == query.lower()), None)
    if exact_match:
        return locations[exact_match], None
    return None, f"Did you mean: {', '.join(suggestions[:3])}?"

# --- Metadata Store ---
metadata_db = {
    'Los Angeles': {'population': '3.8M', 'type': 'city', 'state': 'CA', 'elevation': 71},
    'Las Vegas': {'population': '0.6M', 'type': 'city', 'state': 'NV', 'elevation': 610},
    'San Francisco': {'population': '0.8M', 'type': 'city', 'state': 'CA', 'elevation': 16},
    'New York': {'population': '8.4M', 'type': 'city', 'state': 'NY', 'elevation': 10},
    'Los Alamos': {'population': '12K', 'type': 'town', 'state': 'NM', 'elevation': 2231},
    'Las Cruces': {'population': '100K', 'type': 'city', 'state': 'NM', 'elevation': 1191},
    'Chicago': {'population': '2.7M', 'type': 'city', 'state': 'IL', 'elevation': 179},
    'Miami': {'population': '0.4M', 'type': 'city', 'state': 'FL', 'elevation': 2},
    'Seattle': {'population': '0.7M', 'type': 'city', 'state': 'WA', 'elevation': 54},
    'Denver': {'population': '0.7M', 'type': 'city', 'state': 'CO', 'elevation': 1609}
}

metadata_cache = {}

def get_metadata(location):
    if location in metadata_cache:
        print(f"[Cache Hit] {location}")
        return metadata_cache[location]
    if location in metadata_db:
        print(f"[Cache Miss] Fetching {location} from DB")
        data = metadata_db[location]
        metadata_cache[location] = data
        return data
    return None

# --- Event Handling and Notification System ---
event_queue = queue.Queue()

EVENT_ROUTE_FOUND = 'route_found'
EVENT_TRAFFIC_UPDATE = 'traffic_update'

def log_route_found(data):
    print(f"[Log] Route: {data['path']}, Distance: {data['distance']:.2f} km, Time: {data['time']:.2f} min")

def notify_traffic_update(data):
    print(f"[Notify] Traffic on {data['road']}: {data['delay']} min delay, Speed: {data['speed']} km/h")

subscribers = {
    EVENT_ROUTE_FOUND: [log_route_found],
    EVENT_TRAFFIC_UPDATE: [notify_traffic_update]
}

def publish_event(event_type, data):
    if event_type in subscribers:
        for subscriber in subscribers[event_type]:
            subscriber(data)

def event_worker():
    while True:
        event_type, data = event_queue.get()
        publish_event(event_type, data)
        event_queue.task_done()
        time.sleep(0.1)

threading.Thread(target=event_worker, daemon=True).start()

# --- Graph-Based Road Network Store ---
def haversine(lon1, lat1, lon2, lat2):
    lon1, lat1, lon2, lat2 = map(radians, [lon1, lat1, lon2, lat2])
    dlon = lon2 - lon1
    dlat = lat2 - lat1
    a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
    c = 2 * asin(sqrt(a))
    r = 6371  # Earth's radius in km
    return c * r

G = nx.Graph()
for loc, coord in locations.items():
    G.add_node(loc, pos=coord)

roads = [
    ('Los Angeles', 'Las Vegas', haversine(*locations['Los Angeles'], *locations['Las Vegas']), 1.0, 105),
    ('Las Vegas', 'San Francisco', haversine(*locations['Las Vegas'], *locations['San Francisco']), 1.0, 120),
    ('Los Angeles', 'San Francisco', haversine(*locations['Los Angeles'], *locations['San Francisco']), 1.0, 110),
    ('San Francisco', 'Seattle', haversine(*locations['San Francisco'], *locations['Seattle']), 1.0, 100),
    ('Seattle', 'Denver', haversine(*locations['Seattle'], *locations['Denver']), 1.0, 115),
    ('Denver', 'Chicago', haversine(*locations['Denver'], *locations['Chicago']), 1.0, 105),
    ('Chicago', 'New York', haversine(*locations['Chicago'], *locations['New York']), 1.0, 100),
    ('New York', 'Miami', haversine(*locations['New York'], *locations['Miami']), 1.0, 110),
    ('Las Vegas', 'Denver', haversine(*locations['Las Vegas'], *locations['Denver']), 1.0, 120),
    ('Los Alamos', 'Las Cruces', haversine(*locations['Los Alamos'], *locations['Las Cruces']), 1.0, 90)
]

for road in roads:
    G.add_edge(road[0], road[1], distance=road[2], traffic_factor=road[3], speed_limit=road[4])

def heuristic(u, v):
    u_pos, v_pos = G.nodes[u]['pos'], G.nodes[v]['pos']
    return haversine(u_pos[1], u_pos[0], v_pos[1], v_pos[0])

def find_shortest_path(start, end):
    def weight_func(u, v, d):
        effective_speed = d['speed_limit'] / d['traffic_factor']
        return (d['distance'] / effective_speed) * 60  # Time in minutes
    
    try:
        path = nx.astar_path(G, start, end, heuristic=heuristic, weight=weight_func)
        total_time = sum(weight_func(u, v, G[u][v]) for u, v in zip(path, path[1:]))
        total_distance = sum(G[u][v]['distance'] for u, v in zip(path, path[1:]))
        event_queue.put((EVENT_ROUTE_FOUND, {'path': path, 'distance': total_distance, 'time': total_time}))
        return path, total_distance, total_time
    except nx.NetworkXNoPath:
        return None, float('inf'), float('inf')

def update_traffic(road, new_traffic_factor):
    if G.has_edge(*road):
        G[road[0]][road[1]]['traffic_factor'] = new_traffic_factor
        delay = (new_traffic_factor - 1) * (G[road[0]][road[1]]['distance'] / G[road[0]][road[1]]['speed_limit']) * 60
        speed = G[road[0]][road[1]]['speed_limit'] / new_traffic_factor
        event_queue.put((EVENT_TRAFFIC_UPDATE, {'road': f"{road[0]} to {road[1]}", 'delay': delay, 'speed': speed}))

# --- Flask Routes ---
@app.route('/')
def index():
    return render_template('index.html')

@app.route('/search', methods=['GET'])
def search():
    query = request.args.get('query')
    coords, message = search_location(query)
    if coords:
        return jsonify({'coords': coords, 'message': message})
    return jsonify({'message': message})

@app.route('/metadata', methods=['GET'])
def metadata():
    location = request.args.get('location')
    data = get_metadata(location)
    if data:
        return jsonify(data)
    return jsonify({'error': 'Location not found'})

@app.route('/route', methods=['GET'])
def route():
    start = request.args.get('start')
    end = request.args.get('end')
    path, distance, time = find_shortest_path(start, end)
    if path:
        return jsonify({'path': path, 'distance': distance, 'time': time})
    return jsonify({'error': 'No path found'})

@app.route('/map', methods=['GET'])
def map_view():
    start = request.args.get('start')
    end = request.args.get('end')
    path, _, _ = find_shortest_path(start, end)
    m = folium.Map(location=locations[start], zoom_start=5)
    for loc, coord in locations.items():
        folium.Marker(coord, popup=loc).add_to(m)
    if path:
        path_coords = [locations[loc] for loc in path]
        folium.PolyLine(path_coords, color="blue", weight=2.5, opacity=1).add_to(m)
    return m._repr_html_()

if __name__ == '__main__':
    app.run(debug=True)