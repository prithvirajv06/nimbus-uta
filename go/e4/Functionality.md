function(inputdata){
    if(inputdata.value > 10){
        inputdata.status = "high"
    } else {
        inputdata.status = "low"
    }
    for(var i=0; i< inputdata.items.length; i++){
        var processed = inputdata.items[i].processed
        if(!processed){
            var cupon = inputdata.items[i].cupons
            for(var j=0; j< cupon.length; j++){
                cupon[j].valid = true
                var discountCriteria = cupon[j].discount_criteria
                for (var k=0; k< discountCriteria.length; k++){
                    if(discountCriteria[k].type == "minimum_purchase" && discountCriteria[k].value > 100){
                        cupon[j].valid = false
                    }
                }
            }
            inputdata.items[i].processed = true
        }
    }

return inputdata;
}