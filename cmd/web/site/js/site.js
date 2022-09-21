
let SetUpChanges = function() {
  $(".changeItem").click(function(){
    var addr = $(this).data("addr")
    $.get("/api/v1/change/" + addr, function(data, status) {
      console.log(data)
      $("#change_data").text(JSON.stringify(data, null, 2))
    })
  })
}