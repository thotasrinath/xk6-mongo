import xk6_mongo from 'k6/x/mongo';


/**
 * Index creation on partyTradeDate
 *
 * db.getCollection("testcollection").createIndex({'trade.party.meta.partyTradeDate':1});
 */

const client = xk6_mongo.newClient('mongodb://172.17.0.2:27017');
export default () => {

    var startDate = randomDate(new Date(2000, 0, 1), new Date(2022, 0, 1), 0, 24);

    var endDate = randomDate(startDate, new Date(2022, 0, 1), 0, 24);

    var query = {"trade.party.meta.partyTradeDate": {$gte: startDate, $lte: endDate}};

    console.log("query is :", query);

    var res = client.find("testdb", "testcollection", query, 10);

    console.log(res);
}

function randomDate(start, end, startHour, endHour) {
    var date = new Date(+start + Math.random() * (end - start));
    var hour = startHour + Math.random() * (endHour - startHour) | 0;
    date.setHours(hour);
    return date;
}


