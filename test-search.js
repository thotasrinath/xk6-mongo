import xk6_mongo from 'k6/x/mongo';
import { SharedArray } from 'k6/data';


/**
 * Create index
 *db.getCollection("testcollection").createIndex({ text_sent_150: "text", text_sent_300: "text",text_sent_450:"text" } )
 */

const client = xk6_mongo.newClient('mongodb://172.17.0.2:27017');

const data = new SharedArray('words', function () {
    // All heavy work (opening and processing big files for example) should be done inside here.
    // This way it will happen only once and the result will be shared between all VUs, saving time and memory.
    const f = JSON.parse(open('./words_dictionary.json'));
    return Object.keys(f); // f must be an array
});


function sentenceGenerator(words, size) {
    var sentence = '';
    for (var i = 0; i < size; i++) {
        sentence += words[Math.floor(Math.random() * words.length)] + ' ';
    }

    return sentence;
}


export default () => {


    var query = { $text: { $search: sentenceGenerator(data, 5) } };

    console.log("query is :", query);

    var res = client.find("testdb", "testcollection", query, 10);

    console.log(res);
}



