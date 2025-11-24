const runListObjectsBenchmark = async (fgaClient, modelId, user, relation, type, iterations) => {
  const start = Date.now();
  const options = { authorizationModelId: modelId };
  for (let i = 0; i < iterations; i++) {
    await fgaClient.listObjects({ user, relation, type }, options);
  }
  return (Date.now() - start) / iterations;
};

module.exports = { runListObjectsBenchmark };