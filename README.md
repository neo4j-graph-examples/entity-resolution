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

## Example Use case
This repository covers such a use case of linking similar user accounts for analytics and providing better recommendations. We have taken an example of online movie streaming platform. For ease of understanding, we have taken movies and users datasets. Online streaming platforms are used by families with different profiles for each family members. We are performing Entity Resolution over users’ data to identify similar/same users. We are also performing linking for users which are from same account (or group/family). And later using this linking to provide effective recommendations user.
