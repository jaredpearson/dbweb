# Dreamblade Web

Dreamblade was a miniatures game from Wizards of the Coast that was discontinued in 2007. This fan site provides a catalog of the miniatures, including the stats and abilities of each.

Dreamblade Web is unofficial Fan Content permitted under the Fan Content Policy. Not approved/endorsed by Wizards. Portions of the materials used are property of Wizards of the Coast. ©Wizards of the Coast LLC

## Setup
This site uses MongoDB to manage site specific data. To run an instance of MongoDB within Docker, run the following command.
```
docker run --name mongo -e MONGO_INITDB_ROOT_USERNAME=mongoadmin -e MONGO_INITDB_ROOT_PASSWORD=secret -p 27017:27017 -d mongo:4.0.4
```

## Miniature Data
Data for the miniatures is not included in the source. To get the data, download the Excel data from [BoardGameGeek files](https://boardgamegeek.com/filepage/57443/dreamcatcher-excel) and convert the XLS to CSV. Before starting the application, set the `DATA` environment variable to the path of the CSV file.