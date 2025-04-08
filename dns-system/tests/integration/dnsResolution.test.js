const request = require('supertest');
const app = require('../../backend/app');

describe('DNS Resolution Integration', () => {
  it('should resolve domain', async () => {
    const res = await request(app).get('/api/dns/resolve?domain=google.com');
    expect(res.statusCode).toEqual(200);
    expect(res.body).toHaveProperty('ip');
  });
});