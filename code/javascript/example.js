// npm install --save neo4j-driver
// node example.js
const neo4j = require('neo4j-driver');
const driver = neo4j.driver('neo4j://<HOST>:<BOLTPORT>',
                  neo4j.auth.basic('<USERNAME>', '<PASSWORD>'), 
                  {/* encrypted: 'ENCRYPTION_OFF' */});

const query =
  `
  MATCH (u:User {state: $state} )-[:WATCHED]->(m)-[:HAS]->(g:Genre)
  
  RETURN g.name as genre, count(g) as freq
  ORDER BY freq DESC
  `;

const params = {"state": "Texas"};

const session = driver.session({database:"neo4j"});

session.run(query, params)
  .then((result) => {
    result.records.forEach((record) => {
        console.log(record.get('genre'));
    });
    session.close();
    driver.close();
  })
  .catch((error) => {
    console.error(error);
  });
