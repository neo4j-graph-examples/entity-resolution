# Entity Resolution
What is Entity Resolution (ER)?
Entity Resolution (ER) is the process of disambiguating data to determine if multiple digital records represent the same real-world entity such as a person, organization, place, or other type of object.
For example, say you have information on persons coming from different e-commerce platforms. They may have slightly different contact information, with addresses formatted differently, using different forms/abbreviations of names, etc.
A human may be able to tell if the records actually belong to the same underlying entity but given the number of possible combinations and matching that can be had, there is a need for an intelligent automated approach to doing so, which is where ER systems come into play.

## Use cases
Few of the common and useful entity resolution use cases are below.

### Life Science & Healthcare
Life science and healthcare organizations requires data linking the most. For example, a healthcare organization can implement Entity resolution for consolidation of a patient’s records from a variety of sources, matching data from hospitals and clinics, laboratories, insurance providers and claims and social media profiles to create a unique profile of each patient. This will help providing precise and effective treatment. Similarly, Life science organizations can use ER to connect various entities, research results, input data sets etc. This can facilitate the research & development.
### Insurance and Financial Services
Financial services and Insurance companies often struggle with fragmented and siloed datasets. Because various products\categories maintain their data in different systems and databases. Thus, it is difficult to reconcile a customer's preferences, history, credit ratings etc on a central platform. ER can enable them to perform record linking on different data sets and produce a unified view of customer's state and needs.

### Digital Marketing and content recommendation
Effective marketing and recommendation scheme cannot be produces using distinct data sets or different silos. Records linking, some machine learning and analytics can be very much helpful in producing effective marketing content. Identifying redundant customers is another area in marketing and CRM which needs to be addressed. ER can be mighty effective in such use cases. 

## Demo Use case
This repository covers such a use case of linking similar user accounts for analytics and providing better recommendations. We have taken an example of online movie streaming platform. For ease of understanding, we have taken movies and users datasets. Online streaming platforms are used by families with different profiles for each family members. We are performing Entity Resolution over users’ data to identify similar/same users. We are also performing linking for users which are from same account (or group/family). And later using this linking to provide effective recommendations user. We have taken an example of a dummy online movie streaming platform. For ease of understanding, we have taken only movies and users datasets. Users can have one or more accounts on a movie streaming platform. We are performing Entity Resolution over users’ data to identify similar/same users. We are also performing linking for users which are from same account (or group/family). Later, we are leveraging this linking to provide effective recommendations to individual users.

## Data Model
<img src="documentation\img\model.PNG">

## Preparing the Graph: Loading data and creating Entities and Relationships
In this guide, we will perform below steps:

* Load: Load information from external CSV files and create entities
* Index: Index nodes based on label
* Relate: Establish connections (relationships) between entities
* Test: Perform basic querying with Cypher on loaded data
* ER: Perform Entity Resolution based on similarity and do record linkage
* Recommend: Generate recommendation based on user similarities / preferences
* Additional: Find similar users by preferences

### Notes
In this demonstration, we have used Neo4j APOC (Awesome Procedures on Cypher) and Neo4j GDS (Graph Data Science) libraries few Cypher queries. To execute the Cypher queries with APOC or GDS functions, you will need to add these libraries as plugins to your Neo4j database instance. For more details on APOC and GDS, please refer below links.

[APOC](https://neo4j.com/developer/neo4j-apoc/) 

[GDS](https://neo4j.com/docs/graph-data-science/current/)

To load data from CSV files, you must download the CSV files from below given URL and place them under <NEO4J_HOME>/import/entity-resolution directory. Please make sure you have placed the files before you proceed further.

[Download CSV files for data import](https://github.com/neo4j-graph-examples/enitity_resolution/tree/main/data/csv)

[Guide to import data from CSV](https://neo4j.com/developer/guide-import-csv/)

## Load information from external CSV and create entities

### Load IP Addresses
 ```sh
:auto USING PERIODIC COMMIT 100 LOAD CSV WITH HEADERS FROM "file:///entity-resolution/IPAddress.csv" AS row
CREATE (n:IpAddress)
SET n = row,
n.ipAddressId = toInteger(row.ipAddressId),
n.address = toString(row.address)
```

### Load Users
 ```sh
:auto USING PERIODIC COMMIT 100 LOAD CSV WITH HEADERS FROM "file:///entity-resolution/Users.csv" AS row
CREATE (n:User)
SET n = row,
n.userId = toInteger(row.userId),
n.ipAddressId = toInteger(row.ipAddressId)
```

### Load Genres
 ```sh
:auto USING PERIODIC COMMIT 5 LOAD CSV WITH HEADERS FROM "file:///entity-resolution/Genres.csv" AS row
CREATE (n:Genre)
SET n = row,
n.name = toString(row.name),
n.genreId = toInteger(row.genreId)
```

### Load Movies
 ```sh
:auto USING PERIODIC COMMIT 50 LOAD CSV WITH HEADERS FROM "file:///entity-resolution/Movies.csv" AS row
CREATE (n:Movie)
SET n = row,
n.movieId = toInteger(row.movieId),
n.name = toString(row.name),
n.year = toInteger(row.year),
n.genreId = toInteger(row.genreId)
```

## Index nodes based on label

 ```sh
CREATE INDEX IF NOT EXISTS FOR (g:Genre) ON (g.genreId);
CREATE INDEX IF NOT EXISTS FOR (u:User) ON (u.ipAddressId);
CREATE INDEX IF NOT EXISTS FOR (u1:User) ON (u1.userId);
CREATE INDEX IF NOT EXISTS FOR (m:Movie) ON (m.movieId);
CREATE INDEX IF NOT EXISTS FOR (i:IpAddress) ON (i.ipAddressId);
CREATE INDEX IF NOT EXISTS FOR (u2:User) ON (u2.state);
 ```
 
## Establish connections (relationships) between entities

### Link user and ip addresses
 ```sh
MATCH (u:User)
MATCH(i:IpAddress)
WHERE u.ipAddressId = i.ipAddressId
MERGE (u)-[:USES]->(i)
```
### Link movies and genres
 ```sh
MATCH (g:Genre)
MATCH(m:Movie)
WHERE g.genreId = m.genreId
MERGE (m)-[:HAS]->(g)
```

### Load Watch Events Relationships
 ```sh
LOAD CSV WITH HEADERS FROM "file:///entity-resolution/WatchEvent.csv" AS row
MATCH (u:User {userId: toInteger(row.userId)})
MATCH (m:Movie {movieId: toInteger(row.movieId)})
MERGE (u)-[:WATCHED { timeStamp: row.timestamp }]->(m)
```

## Perform basic querying with Cypher on loaded data

### Query users who have watched movie "The Boss Baby: Family Business"

```sh
MATCH (u:User)-[w:WATCHED]->(m:Movie {name: "The Boss Baby: Family Business"}) RETURN u,w,m LIMIT 5
Show users from "New York" and movies watched by them

MATCH (u:User {state: "New York"} )-[w:WATCHED]->(m:Movie)  RETURN u,w,m LIMIT 50
Show trending genres in Texas

MATCH (u:User {state: "Texas"} )-[w:WATCHED]->(m:Movie)-->(g:Genre)
return g.name, count(g) order by count(g) desc
```


## Perform Entity Resolution based on similarity and do record linkage

### Users who have similar names
These are users who have same/similar names but different (redundant) profiles due to typos or abbreviations used for some instances. We are using the Jaro Winkler Distance algorithm from the Neo4j APOC library.

References

[Jaroâ€“Winkler distance](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)

[apoc.text.jaroWinklerDistance](https://neo4j.com/labs/apoc/4.1/overview/apoc.text/apoc.text.jaroWinklerDistance/)

```sh
MATCH (a:User)
MATCH (b:User)
WHERE a.firstName + a.lastName <> b.firstName + b.lastName
WITH a, b, a.firstName + a.lastName AS norm1, b.firstName + b.lastName AS norm2
WITH
toInteger(apoc.text.jaroWinklerDistance(norm1, norm2) * 100) AS nameSimilarity,
toInteger(apoc.text.jaroWinklerDistance(a.email, b.email) * 100) AS emailSimilarity,
toInteger(apoc.text.jaroWinklerDistance(a.phone, b.phone) * 100) AS phoneSimilarity, a, b
WITH a, b, toInteger((nameSimilarity + emailSimilarity + phoneSimilarity)/3) as similarity WHERE similarity >= 90
RETURN a.firstName + a.lastName AS p1, b.firstName + b.lastName AS p2, a.email, b.email,  similarity
```

### Users belonging to same family
Users who have similar last names and live in same state, and use same IP address, that means they are either same users with redundant profile or belong to the same family

```sh
MATCH (a:User)-[:USES]->(:IpAddress)<-[:USES]-(b:User)
WHERE a.lastName =  b.lastName AND a.state = b.state AND a.country = b.country
WITH a.lastName as familyName, collect(distinct b) as familyMembers, count(distinct b) as totalMembers
UNWIND  familyMembers as member
RETURN familyName, totalMembers,  member.firstName + ' '  + member.lastName  AS memberName
```

### Record Linkage: Create Family Nodes for each family and connect members. This is how we link the similar users and family members using a common Family node
```sh
MATCH (a:User)-[:USES]->(:IpAddress)<-[:USES]-(b:User)
WHERE a.lastName =  b.lastName AND a.state = b.state AND a.country = b.country
WITH a.lastName as familyName, collect(distinct b) as familyMembers, count(distinct b) as totalMembers
MERGE (a:Family {name: familyName})
WITH a,familyMembers
UNWIND  familyMembers as member
MERGE (member)-[:BELONGS_TO]->(a)
Check how may families are created
MATCH (f:Family)<-[:BELONGS_TO]-(u:User) RETURN f, u LIMIT 200
```

### Check how may families are created
```sh
MATCH (f:Family)<-[:BELONGS_TO]-(u:User) RETURN f, u LIMIT 200
```

## Generate recommendation based on user similarities / preferences

Providing recommendation to the member based on his/her account/family members history. Get preferred genres by other account members and suggest top 5 movies from most watched genres.

```sh
MATCH (user:User {firstName: "Vilma", lastName: "De Mars"})
MATCH (user)-[:BELONGS_TO]->(f)<-[:BELONGS_TO]-(otherMember)
MATCH (otherMember)-[:WATCHED]->(m:Movie)-[:HAS]->(g:Genre)<-[:HAS]-(m2:Movie)
WITH g.name as genre, count(distinct m2) as totalMovies, collect(m2.name) as movies
RETURN genre, totalMovies, movies[0..5] as topFiveMovies ORDER BY totalMovies DESC LIMIT 50
```

## Find similar users by preferences

We can also find similarity in records using properties and/or connections. We will perform couple of examples to demonstrate the same

### Example 1: Find users based on their movie watching preferences using Node Similarity algorithm
[Node Similarity](https://neo4j.com/docs/graph-data-science/current/algorithms/node-similarity/)

#### Step 1: For this, we will first create an in-memory graph with node and relationship specification to perform matching.
```sh
CALL gds.graph.create(
    'similarityGraph',
    ['User', 'Movie'],
    {
        WATCHED: {
            type: 'WATCHED'
        }
    }
);
```
#### Step 2: Perform memory estimate for the matching to execute
```sh
CALL gds.nodeSimilarity.write.estimate('similarityGraph', {
  writeRelationshipType: 'SIMILAR',
  writeProperty: 'score'
})
YIELD nodeCount, relationshipCount, bytesMin, bytesMax, requiredMemory
```
#### Step 3: Execute algorithm and show results.
```sh
CALL gds.nodeSimilarity.stream('similarityGraph')
YIELD node1, node2, similarity
WITH gds.util.asNode(node1) AS Person1, gds.util.asNode(node2) AS Person2, similarity
RETURN
Person1.firstName + ' ' +  Person1.lastName as p1,
Person2.firstName  + ' ' +   Person2.lastName as p2, similarity ORDER BY similarity DESC
```


### Example 2: Find similar users by genre preference using Pearson similarity function
[Peason Similarity - Neo4j GDS](https://neo4j.com/docs/graph-data-science/current/alpha-algorithms/pearson/)

[Pearson correlation coefficient](https://en.wikipedia.org/wiki/Pearson_correlation_coefficient)

Here we are finding the users who have similar preferences as user "Lanette Laughtisse". We are comparing the similarities based on the movies they have watched from similar genre. We can use this information to provide recommendations.
```sh
MATCH (p1:User {firstName:"Lanette", lastName:"Laughtisse"} )-[:WATCHED]->(m:Movie)
MATCH (m)-[:HAS]->(g1:Genre)
WITH p1, g1, count(m) as movieCount1
WITH p1, gds.alpha.similarity.asVector(g1, movieCount1) AS p1Vector
MATCH (p2:User)-[:WATCHED]->(m2:Movie)
MATCH (m2)-[:HAS]->(g1:Genre) WHERE p2 <> p1
WITH p1, g1, p1Vector, p2, count(m2) as movieCount2
WITH p1, p2, p1Vector, gds.alpha.similarity.asVector(g1, movieCount2) AS p2Vector
WHERE size(apoc.coll.intersection([v in p1Vector | v.category], [v in p2Vector | v.category])) > 3
WITH
p1.firstName + ' '  + p1.lastName  AS currentUser,
p2.firstName + ' ' + p2.lastName  AS similarUser,
gds.alpha.similarity.pearson(p1Vector, p2Vector, {vectorType: "maps"}) AS similarity
WHERE similarity > 0.9
RETURN currentUser,similarUser, similarity
       ORDER BY similarity DESC
LIMIT 100
```

## References

* [Developer resources](https://neo4j.com/developer/)
* [Neo4j Cypher Manual](https://neo4j.com/docs/cypher-manual/current/)
* [Entity Resolution in Neo4j reference](https://neo4j.com/developer-blog/exploring-supervised-entity-resolution-in-neo4j/)
