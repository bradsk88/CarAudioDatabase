import {Component, OnInit} from '@angular/core';
import * as Highcharts from 'highcharts';
import {XAxisOptions} from "highcharts";

@Component({
  selector: 'app-graph',
  templateUrl: './graph.component.html',
  styleUrls: ['./graph.component.scss']
})
export class GraphComponent implements OnInit {

  Highcharts: typeof Highcharts = Highcharts;
  chartOptions: Highcharts.Options = {
    series: [{
      data: [[1, 1], [2, 2], [3, 3]],
      type: 'line',
    }],
    xAxis: {
      type: "logarithmic",
    }
  };

  results: any = [
    {
      name: "Frequency Response",
      series: [
        {name: 1, value: 100},
        {name: 2, value: 70},
      ]
    },
  ];

  constructor() {
  }

  ngOnInit(): void {
    fetch('http://190.92.153.141:8080/get').
    then((v: Response) => v.json()).
    then((v: {data: any[]}) => this.chartOptions = {
      series: [{
        data: v.data.map((v: any) => ([v.Frequency as number, v.Amplitude as number])),
        // data: [[1, 1], [2, 2], [3, 3]],
        type: 'line',
      }],
    });
  }

}
