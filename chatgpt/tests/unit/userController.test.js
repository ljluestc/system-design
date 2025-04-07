const { register } = require('../../backend/controllers/userController');
const User = require('../../backend/models/User');
const bcrypt = require('bcrypt');

jest.mock('../../backend/models/User');
jest.mock('bcrypt');

describe('User Controller', () => {
  it('should register a user', async () => {
    const req = { body: { username: 'test', password: 'pass' } };
    const res = { status: jest.fn().mockReturnThis(), json: jest.fn() };
    bcrypt.hash.mockResolvedValue('hashed');
    User.prototype.save = jest.fn().mockResolvedValue();
    await register(req, res);
    expect(res.status).toHaveBeenCalledWith(201);
  });
});