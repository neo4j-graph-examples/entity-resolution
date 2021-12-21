# Entity Resolution
What is Entity Resolution (ER)?
Entity Resolution (ER) is the process of disambiguating data to determine if multiple digital records represent the same real-world entity such as a person, organization, place, or other type of object.
For example, say you have information on persons coming from different e-commerce platforms. They may have slightly different contact information, with addresses formatted differently, using different forms\abbreviations of names, etc.
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
This repository covers such a use case of linking similar user accounts for analytics and providing better recommendations. We have taken an example of a dummy online movie streaming platform. For ease of understanding, we have taken only movies and users datasets. Users can have one or more accounts on a movie streaming platform. We are performing Entity Resolution over users’ data to identify similar/same users. We are also performing linking for users which are from same account (or group/family). Later, we are leveraging this linking to provide effective recommendations to individual users.

## Data Model
<img src="documentation\img\model.PNG">

## Preparing the Graph: Loading data and creating Entities and Relationships
In this guide, we will perform below steps:

* Load: Load nodes and relationship information from external CSV files and create entities
* Relate: Establish more connections (relationships) between entities
* Test: Perform basic querying with Cypher on loaded data
* ER: Perform Entity Resolution based on similarity and do record linkage
* Recommend: Generate recommendation based on user similarities / preferences
* Additional: Try couple of preference based similarities and recommendation examples

### Notes
In this demonstration, we have used Neo4j APOC (Awesome Procedures on Cypher) and Neo4j GDS (Graph Data Science) libraries few Cypher queries. To execute the Cypher queries with APOC or GDS functions, you will need to add these libraries as plugins to your Neo4j database instance. For more details on APOC and GDS, please refer below links.

[APOC](https://neo4j.com/developer/neo4j-apoc/) 

[GDS](https://neo4j.com/docs/graph-data-science/current/)

## Load nodes and relationship information from external CSV files and create entities

### Load Users, Ip Addresses and connect Users with IP Addresses
 ```sh
// Constraints
CREATE CONSTRAINT user_id IF NOT EXISTS FOR (u:User) REQUIRE u.userId IS UNIQUE;
CREATE CONSTRAINT ip_address IF NOT EXISTS FOR (i:IpAddress) REQUIRE i.address IS UNIQUE;
 
// Data load
LOAD CSV WITH HEADERS FROM "https://gist.githubusercontent.com/chintan196/6b33019341bdcb6ed4d712cc94b84fc6/raw/2513454dd72b70d3122fd0a15777fc9842bbba89/Users.csv" AS row
MERGE (u:User { userId: toInteger(row.userId) })
ON CREATE SET 
u.firstName= row.firstName,
u.lastName= row.lastName,
u.gender= row.gender,
u.email= row.email,
u.phone= row.phone,
u.state= row.state,
u.country= row.country
WITH u, row
MERGE (ip:IpAddress { address: row.ipAddress })
MERGE (u)-[:USES]->(ip)
RETURN u, ip
```

### Load Movies, Genres and link them
 ```sh
// Constraints
CREATE CONSTRAINT genre_name IF NOT EXISTS FOR (g:Genre) REQUIRE g.name IS UNIQUE;
CREATE CONSTRAINT movie_id IF NOT EXISTS FOR (m:Movie) REQUIRE m.movieId IS UNIQUE;
CREATE CONSTRAINT movie_name IF NOT EXISTS FOR (m:Movie) REQUIRE m.name IS UNIQUE;

//Load Data
:auto USING PERIODIC COMMIT 500
LOAD CSV WITH HEADERS FROM 
"https://gist.githubusercontent.com/chintan196/6b33019341bdcb6ed4d712cc94b84fc6/raw/2513454dd72b70d3122fd0a15777fc9842bbba89/Movies.csv" AS row
MERGE ( m:Movie { movieId: toInteger(row.movieId) })
ON CREATE SET 
m.name= row.name,
m.year= toInteger(row.year)
WITH m, row
MERGE (g:Genre { name: row.genre } )
MERGE (m)-[:HAS]->(g) RETURN m, g;
```

## Create connections (relationships) between entities

### Load Watch Events Relationships - Execute this after loading user and movies
 ```sh
LOAD CSV WITH HEADERS FROM "https://gist.githubusercontent.com/chintan196/6b33019341bdcb6ed4d712cc94b84fc6/raw/2513454dd72b70d3122fd0a15777fc9842bbba89/WatchEvent.csv" AS row
MATCH (u:User {userId: toInteger(row.userId)})
MATCH (m:Movie {movieId: toInteger(row.movieId)})  
MERGE (u)-[w:WATCHED]->(m) ON CREATE SET w.watchCount = toInteger(row.watchCount)
RETURN u, m;
```

## Perform basic querying on loaded data

```sh
//Query users who have watched movie "The Boss Baby: Family Business"
MATCH (u:User)-->(m:Movie {name: "The Boss Baby: Family Business"}) RETURN u,m LIMIT 5

//Show users from "New York" and movies watched by them
MATCH (u:User {state: "New York"} )-[:WATCHED]->(m)  RETURN u, m LIMIT 50

//Show trending genres in Texas
MATCH (u:User {state: "Texas"} )-[:WATCHED]->(m)-[:HAS]->(g)
return g.name, count(g) order by count(g) desc
```


## Perform Entity Resolution based on similarity and do record linkage

### Users who have similar names
These are users who have same/similar names but different (redundant) profiles due to typos or abbreviations used for some instances. We are using the Jaro Winkler Distance algorithm from the Neo4j APOC library.

References

[Jaro Winkler distance](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)

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
MATCH (a:User)-->(:IpAddress)<--(b:User)
WHERE a.lastName =  b.lastName AND a.state = b.state AND a.country = b.country
WITH a.lastName as familyName, collect(distinct b.firstName + ' '  + b.lastName) as members, count(distinct b) as memberCount
RETURN familyName, memberCount, members
```

### Record Linkage: Create Family Nodes for each family and connect members. This is how we link the similar users and family members using a common Family node
```sh
MATCH (a:User)-->(:IpAddress)<--(b:User)
WHERE a.lastName =  b.lastName AND a.state = b.state AND a.country = b.country
WITH a.lastName as familyName, collect(distinct b) as familyMembers, count(distinct b) as totalMembers
MERGE (a:Family {name: familyName})
WITH a,familyMembers
UNWIND  familyMembers as member
MERGE (member)-[:BELONGS_TO]->(a)
RETURN a, member
```

### Check how may families are created
```sh
MATCH (f:Family)<--(u:User) RETURN f, u LIMIT 200
```

## Generate recommendation based on user similarities / preferences

Providing recommendation to the member based on his/her account/family members history. Get preferred genres by other account members and suggest top 5 movies from most watched genres.

```sh
MATCH (user:User {firstName: "Vilma", lastName: "De Mars"})
MATCH (user)-[:BELONGS_TO]->(f)<-[:BELONGS_TO]-(otherMember)
MATCH (otherMember)-[:WATCHED]->(m1)-[:HAS]->(g:Genre)<-[:HAS]-(m2)
WITH g.name as genre, count(distinct m2) as totalMovies, collect(m2.name) as movies
RETURN genre, totalMovies, movies[0..5] as topFiveMovies ORDER BY totalMovies DESC LIMIT 50 
```

## Using Neo4j Node Similarity Algorigthm to find similar users and get recommendations

[Node Similarity](https://neo4j.com/docs/graph-data-science/current/algorithms/node-similarity/)

#### Step 1: For this, we will first create an in-memory graph with node and relationship specification to perform matching.
```sh
CALL gds.graph.create(
    'similarityGraph',
    ['User', 'Movie'],
    {
        WATCHED: {
            type: 'WATCHED',
            properties: {
                strength: {
                    property: 'watchCount',
                    defaultValue: 1
                }
            }
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
#### Step 4: Get recommendations for a user based on similarity. For a user, fetch recommendations based on other similar users' preferences
```sh
CALL gds.nodeSimilarity.stream('similarityGraph')
YIELD node1, node2, similarity
WITH gds.util.asNode(node1) AS Person1, gds.util.asNode(node2) AS Person2, similarity
WHERE Person1.firstName = 'Paulie' AND Person1.lastName = 'Imesson'
MATCH (Person2)-[w:WATCHED]->(m) WHERE NOT exists((Person1)-->(m))
WITH  DISTINCT m as movies, SUM(w.watchCount) as watchCount
RETURN movies order by watchCount
```

## Using Pearson Similarity Algorigthm to find similar users based on Genre preference and get recommendations

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
Get recommendations for a user using similar order users' preferenes by fetching similar users using Pearson Similarity function

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
p1 AS currentUser,
p2 AS similarUser,
gds.alpha.similarity.pearson(p1Vector, p2Vector, {vectorType: "maps"}) AS similarity
WHERE similarity > 0.9
MATCH (similarUser)-[w:WATCHED]->(m) 
WITH  DISTINCT m as movies, SUM(w.watchCount) as watchCount
RETURN movies order by watchCount
```

## References

* [Developer resources](https://neo4j.com/developer/)
* [Neo4j Cypher Manual](https://neo4j.com/docs/cypher-manual/current/)
* [Entity Resolution in Neo4j reference](https://neo4j.com/developer-blog/exploring-supervised-entity-resolution-in-neo4j/)
