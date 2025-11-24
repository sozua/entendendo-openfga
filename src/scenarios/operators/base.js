const { writeBatches } = require('../../utils/seeder');

const generator = async (fgaClient) => {
  const body = {
    schema_version: '1.1',
    type_definitions: [
      { type: 'user' },
      {
        type: 'document',
        relations: { 
          viewer: { this: {} } 
        },
        metadata: { 
          relations: { 
            viewer: { directly_related_user_types: [{ type: 'user' }] } 
          } 
        }
      }
    ]
  };
  const resp = await fgaClient.writeAuthorizationModel(body);
  return resp.authorization_model_id;
};

const seeder = async (fgaClient, count) => {
  const tuples = [];
  for(let i=1; i<=count; i++) {
    tuples.push({ user: 'user:benchmark_user', relation: 'viewer', object: `document:${i}` });
  }
  await writeBatches(fgaClient, tuples);
};

module.exports = {
  id: 'base',
  name: 'Base (Direct Assignment)',
  generator,
  seeder,
  checkTuple: { user: 'user:benchmark_user', relation: 'viewer', object: 'document:1' }
};