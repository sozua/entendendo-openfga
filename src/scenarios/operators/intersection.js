const { writeBatches } = require('../../utils/seeder');

const generator = async (fgaClient) => {
  const body = {
    schema_version: '1.1',
    type_definitions: [
      { type: 'user' },
      {
        type: 'document',
        relations: { 
          a: { this: {} },
          b: { this: {} },
          viewer: { intersection: { child: [{ computedUserset: { relation: 'a' } }, { computedUserset: { relation: 'b' } }] } }
        },
        metadata: { 
          relations: { 
            a: { directly_related_user_types: [{ type: 'user' }] },
            b: { directly_related_user_types: [{ type: 'user' }] }
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
    tuples.push({ user: 'user:benchmark_user', relation: 'a', object: `document:${i}` });
    tuples.push({ user: 'user:benchmark_user', relation: 'b', object: `document:${i}` });
  }
  await writeBatches(fgaClient, tuples);
};

module.exports = {
  id: 'and',
  name: 'Intersection (A AND B)',
  generator,
  seeder,
  checkTuple: { user: 'user:benchmark_user', relation: 'viewer', object: 'document:1' }
};