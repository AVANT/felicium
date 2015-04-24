{
    "media": [
    '{{repeat(26)}}',
         {
             "altText": "{{lorem(1, sentences)}}",
             "Title": "User {{index}} image.",
             "seedData":{
                 "url": function(idx){return "http://placekitten.com/200/" + (200+idx);},
                 "belongsTo": {
                     "user": "{{index}}"
                 },
                 "upoadedBy": "{{index}}"
             }
         }  
	],
    "users":[
        '{{repeat(20}}'
    ]
}
