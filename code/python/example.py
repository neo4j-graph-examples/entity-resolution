# pip3 install neo4j
# python3 example.py

from neo4j import GraphDatabase, basic_auth

driver = GraphDatabase.driver(
  "neo4j://<HOST>:<BOLTPORT>",
  auth=basic_auth("<USERNAME>", "<PASSWORD>"))

cypher_query = '''
MATCH (u:User {state: $state} )-[:WATCHED]->(m)-[:HAS]->(g:Genre)

RETURN g.name as genre, count(g) as freq
ORDER BY freq DESC
'''

with driver.session(database="neo4j") as session:
  results = session.read_transaction(
    lambda tx: tx.run(cypher_query,
                      state="Texas").data())
  for record in results:
    print(record['genre'])

driver.close()
