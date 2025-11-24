const { BATCH_SIZE, CONCURRENCY } = require('../config');

const writeBatches = async (fgaClient, allTuples) => {
  const queue = [];
  for (let i = 0; i < allTuples.length; i += BATCH_SIZE) {
    queue.push(allTuples.slice(i, i + BATCH_SIZE));
  }

  async function worker() {
    while (queue.length > 0) {
      const batch = queue.shift();
      try {
        await fgaClient.writeTuples(batch);
      } catch (e) {
        console.error('Error seeding batch:', e);
      }
    }
  }

  const workers = [];
  for (let i = 0; i < CONCURRENCY; i++) workers.push(worker());
  await Promise.all(workers);
};

module.exports = { writeBatches };