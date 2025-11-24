const { writeBatches } = require('../../utils/seeder');

const generator = async (fgaClient) => {
  const typeDefinitions = [{ type: 'user' }];
  for (let i = 0; i < 80; i++) {
    typeDefinitions.push({
      type: `type_${String(i).padStart(2, '0')}`,
      relations: { viewer: { this: {} } },
      metadata: {
        relations: {
          viewer: { directly_related_user_types: [{ type: 'user' }] },
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
  const tuples = [];
  for (let i = 0; i < count; i++) {
    const typeIdx = i % 80;
    tuples.push({
      user: 'user:benchmark_user',
      relation: 'viewer',
      object: `type_${String(typeIdx).padStart(2, '0')}:${i}`,
    });
  }
  await writeBatches(fgaClient, tuples);
};

module.exports = {
  id: 'wide',
  name: 'Modelo Amplo (Muitos Tipos, Baixa Profundidade)',
  generator,
  seeder,
  checkTuple: { user: 'user:benchmark_user', relation: 'viewer', object: 'type_00:1' }
};