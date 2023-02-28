import xk6_mongo from 'k6/x/mongo';


/**
 * Create index
 *db.getCollection("testcollection").createIndex({ text_sent_150: "text", text_sent_300: "text",text_sent_450:"text" } )
 */

const client = xk6_mongo.newClient('mongodb://172.17.0.2:27017');
export default () => {


    var query = { $text: { $search: "acetabularia" } };

    console.log("query is :", query);

    var res = client.find("testdb", "testcollection", query, 10);

    console.log(res);
}



