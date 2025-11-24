const { writeBatches } = require('../../utils/seeder');

const generator = async (fgaClient) => {
  const typeDefinitions = [{ type: 'user' }];
  const DEPTH = 15;

  for (let i = 0; i < DEPTH; i++) {
    const currentType = `node_${String(i).padStart(2, '0')}`;
    const nextType =
      i === DEPTH - 1 ? 'user' : `node_${String(i + 1).padStart(2, '0')}`;

    typeDefinitions.push({
      type: currentType,
      relations: { next: { this: {} } },
      metadata: {
        relations: {
          next: {
            directly_related_user_types: [{ type: nextType }, { type: 'user' }],
          },
        },
      },
    });
  }

  const resp = await fgaClient.writeAuthorizationModel({
    schema_version: '1.1',
    type_definitions: typeDefinitions,
  });
  return resp.authorization_model_id;
};

const seeder = async (fgaClient, count) => {
  const DEPTH = 15;
  const tuples = [];

  for (let i = 0; i < DEPTH; i++) {
    const currentObj = `node_${String(i).padStart(2, '0')}:1`;
    const nextObj =
      i === DEPTH - 1
        ? 'user:benchmark_user'
        : `node_${String(i + 1).padStart(2, '0')}:1`;
    tuples.push({ user: nextObj, relation: 'next', object: currentObj });
  }

  for (let i = 0; i < count; i++) {
    const level = i % DEPTH;
    const currentObj = `node_${String(level).padStart(2, '0')}:${i + 100}`;
    tuples.push({ user: 'user:random', relation: 'next', object: currentObj });
  }

  await writeBatches(fgaClient, tuples);
};

module.exports = {
  id: 'deep',
  name: 'Modelo Profundo (Poucos Tipos, Alta Profundidade 15)',
  generator,
  seeder,
  checkTuple: { user: 'user:benchmark_user', relation: 'next', object: 'node_00:1' }
};