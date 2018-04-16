function exportSummary()
{
  // This should probably be async....
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.open( "POST", "http://localhost:8080/export_inventory_summary", false );
  xmlHttp.send( null );
  console.log(xmlHttp.status)
  if(xmlHttp.status == 204) {
    document.getElementById("export_summary").innerHTML = "<button type=\"button\" class=\"oops-button\" id=\"export_summary\" disabled>Export Summary to CSV</button>";
    document.getElementById("export_message").innerHTML = "Export successful! CSV's have been generated!";
  }
}

function getStateTrackerState()
{
  // This should probably be async....
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.open( "GET", "http://localhost:8080/state", false );
  xmlHttp.send( null );
  return JSON.parse(xmlHttp.responseText);
}

function drawLine(ctx, startX, startY, endX, endY, color){
    ctx.save();
    ctx.strokeStyle = color;
    ctx.beginPath();
    ctx.moveTo(startX,startY);
    ctx.lineTo(endX,endY);
    ctx.stroke();
    ctx.restore();
}

function drawBar(ctx, upperLeftCornerX, upperLeftCornerY, width, height, color){
    ctx.save();
    ctx.fillStyle=color;
    ctx.fillRect(upperLeftCornerX, upperLeftCornerY, width, height);
    ctx.restore();
}

function millisToDateTime(millis) {
  return new Date(millis);
}

function sortChanges(a, b) {
  if (a.created_at < b.created_at)
    return -1;
  if (a.created_at > b.created_at)
    return 1;
  return 0;
}

// options:
/*
padding, colors, gridColor,  gridYLines,
*/
var ChangeChart = function(options, canvas, changes) {
    this.options = options;
    this.canvas = canvas;
    this.changes = changes;
    this.changes.sort(sortChanges);
    this.ctx = this.canvas.getContext("2d");
    this.colors = options.colors;

    this.draw = function(){
        var maxValue = 0;
        var minValue = 0;
        var millisPerDay = 8.64e+7;

        var numChanges = changes.length;
        var minMillis = this.changes[0].created_at;
        var maxMillis = this.changes[numChanges - 1].created_at;

        // Figure out date buckets
        var startMidnight = millisToDateTime(minMillis);
        console.log(startMidnight.toLocaleString());
        startMidnight.setHours(24, 0, 0, 0);
        var endMidnight = millisToDateTime(maxMillis);
        console.log(endMidnight.toLocaleString());
        endMidnight.setHours(24, 0, 0, 0);
        var maxMillis = endMidnight.getTime()
        var totalMillis = maxMillis - minMillis;
        var numDays = Math.ceil(1.0 * totalMillis / millisPerDay);
        console.log(startMidnight.toLocaleString(), endMidnight.toLocaleString());

        var canvasActualHeight = this.canvas.height - this.options.padding * 2;
        var canvasActualWidth = this.canvas.width - this.options.padding * 2;
        // bars have width of 15 mins
        var barWidth = (canvasActualWidth) / (numDays * 24 * 15);

        // Get max and min value of the sum.
        var sum = 0;
        for (var cc in this.changes) {
            var mode = this.changes[cc].mode
            var quantity = this.changes[cc].quantity
            sum = sum + mode * quantity;
            maxValue = Math.max(maxValue, sum);
            minValue = Math.min(minValue, sum);
        }

        // Determine Y span
        var gridYSpan = maxValue;
        if (minValue <= 0) {
          gridYSpan = maxValue - minValue;
        }

        var gridYScale = Math.round(1.0 * gridYSpan / this.options.gridYLines);
        var barMinValue = canvasActualHeight * minValue / gridYSpan;
        var barMaxValue = canvasActualHeight * maxValue / gridYSpan;

        // label day buckets
        var dateMillisValue = startMidnight.getTime();
        while(dateMillisValue <= maxMillis) {
          var changeTimePct = 1.0 * (dateMillisValue - minMillis) / totalMillis;
          var date = new Date(dateMillisValue);
          var gridXValue = changeTimePct * canvasActualWidth;
          var gridYValue = canvasActualHeight + this.options.padding;
          drawLine(
              this.ctx,
              gridXValue,
              gridYValue,
              gridXValue,
              gridYValue - 20,
              this.options.gridColor
          );
          //writing grid markers
          this.ctx.save();
          this.ctx.fillStyle = this.options.gridColor;
          var labelDate = new Date(dateMillisValue - millisPerDay);
          this.ctx.fillText(labelDate.toDateString(), gridXValue-100, gridYValue + this.options.padding / 2);
          this.ctx.restore();

          dateMillisValue+=millisPerDay;
        }

        //drawing the positive grid lines
        var gridValue = 0;
        while (gridValue <= maxValue) {
            var gridY = canvasActualHeight * (1 - (gridValue)/gridYSpan) + barMinValue + this.options.padding;
            drawLine(
                this.ctx,
                0,
                gridY,
                this.canvas.width,
                gridY,
                this.options.gridColor
            );

            //writing grid markers
            this.ctx.save();
            this.ctx.fillStyle = this.options.gridColor;
            this.ctx.fillText(gridValue, 12, gridY - 2);
            this.ctx.restore();

            gridValue+=gridYScale;
        }

        //drawing the positive grid lines
        var gridValue = 0;
        while (gridValue >= minValue) {
            var gridY = canvasActualHeight * (1 - (gridValue)/gridYSpan) + barMinValue + this.options.padding;
            drawLine(
                this.ctx,
                0,
                gridY,
                this.canvas.width,
                gridY,
                this.options.gridColor
            );

            //writing grid markers
            this.ctx.save();
            this.ctx.fillStyle = this.options.gridColor;
            this.ctx.fillText(gridValue, 12, gridY - 2);
            this.ctx.restore();

            gridValue-=gridYScale;
        }

        //drawing the bars
        var sum = 0;
        var lastX = 0;
        var lastY = 0;
        for (cc in this.changes) {
            var mode = this.changes[cc].mode
            var quantity = this.changes[cc].quantity
            sum = sum + mode * quantity
            var barHeight = Math.round( canvasActualHeight * sum / gridYSpan);
            var changeTimePct = 1.0 * (this.changes[cc].created_at - minMillis) / totalMillis;
            var barXIndex = changeTimePct * canvasActualWidth;
            var rectX = this.options.padding + barXIndex;
            var rectY = this.canvas.height - barHeight + barMinValue - this.options.padding;
            var ptX = rectX;
            var ptY = rectY;
            if(sum < 0) {
              rectY = this.canvas.height + barMinValue - this.options.padding;
            }

            drawBar(
                this.ctx,
                rectX,
                rectY,
                barWidth,
                Math.abs(barHeight),
                this.colors[1]
            );

            if(cc > 0) {
              drawLine(
                  this.ctx,
                  lastX,
                  lastY,
                  ptX,
                  ptY,
                  this.colors[0]
              );
            }
            lastX = ptX;
            lastY = ptY;
        }

    }
}

function drawSummaryCharts() {
  var fullState = getStateTrackerState();
  for (ss in fullState.inventory_summary) {
    var chartCanvas = document.getElementById(fullState.inventory_summary[ss].item_code);
    // Make it visually fill the positioned parent
    chartCanvas.style.width ='100%';
    // ...then set the internal size to match
    chartCanvas.width  = chartCanvas.offsetWidth;
    chartCanvas.height = 300;
    var ssChangeChart = new ChangeChart(
      {
        padding:30,
        gridYLines:10,
        gridColor:"#84817c",
        colors:["#eb9743","#67b6c7"]
      },
      chartCanvas,
      fullState.inventory_summary[ss].inventory_changes);
    ssChangeChart.draw();
  }
}

document.onload = drawSummaryCharts();
