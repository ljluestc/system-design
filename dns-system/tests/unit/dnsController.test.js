const { resolveDomain } = require('../../backend/controllers/dnsController');
const { getFromCache, setInCache } = require('../../backend/services/cacheService');

jest.mock('../../backend/services/cacheService');

describe('dnsController', () => {
  it('should resolve domain from cache', async () => {
    getFromCache.mockResolvedValue('192.168.1.1');
    const req = { query: { domain: 'example.com' } };
    const res = { json: jest.fn() };
    await resolveDomain(req, res);
    expect(res.json).toHaveBeenCalledWith({ ip: '192.168.1.1' });
  });
});