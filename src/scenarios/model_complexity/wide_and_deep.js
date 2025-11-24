const { writeBatches } = require('../../utils/seeder');

const generator = async (fgaClient) => {
  const typeDefinitions = [{ type: 'user' }];
  for (let i = 0; i < 65; i++) {
    typeDefinitions.push({
      type: `junk_${i}`,
      relations: { viewer: { this: {} } },
      metadata: {
        relations: {
          viewer: { directly_related_user_types: [{ type: 'user' }] },
        },
      },
    });
  }

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
          next: { directly_related_user_types: [{ type: nextType }] },
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
    tuples.push({
      user: 'user:random',
      relation: 'viewer',
      object: `junk_${i % 65}:${i}`,
    });
  }

  await writeBatches(fgaClient, tuples);
};

module.exports = {
  id: 'wide_deep',
  name: 'Modelo Amplo e Profundo (Muitos Tipos + Alta Profundidade)',
  generator,
  seeder,
  checkTuple: { user: 'user:benchmark_user', relation: 'next', object: 'node_00:1' }
};