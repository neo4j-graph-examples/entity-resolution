# pip3 install neo4j
# python3 example.py

from neo4j import GraphDatabase, basic_auth

cypher_query = '''
MATCH (u:User {state: $state} )-[:WATCHED]->(m)-[:HAS]->(g:Genre)

RETURN g.name as genre, count(g) as freq
ORDER BY freq DESC
'''

with GraphDatabase.driver(
    "neo4j://<HOST>:<BOLTPORT>",
    auth=("<USERNAME>", "<PASSWORD>")
) as driver:
    result = driver.execute_query(
        cypher_query,
        state="Texas",
        database_="neo4j")
    for record in result.records:
        print(record['genre'])
