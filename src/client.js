const { OpenFgaClient } = require('@openfga/sdk');
const { OPENFGA_URL } = require('./config');

const getClient = () => {
  return new OpenFgaClient({
    apiUrl: OPENFGA_URL,
  });
};

const waitForOpenFGA = async () => {
  for (let i = 0; i < 60; i++) {
    try {
      const resp = await fetch(`${OPENFGA_URL}/healthz`);
      if (resp.ok) return;
    } catch (e) {}
    await new Promise((r) => setTimeout(r, 1000));
  }
  throw new Error('OpenFGA failed to start');
};

module.exports = {
  getClient,
  waitForOpenFGA
};