import {Component, Inject, OnInit} from '@angular/core';
import * as Highcharts from 'highcharts';
import {DOCUMENT} from "@angular/common";

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
      tickInterval: 0.2,
      min: 15,
      max: 22000,
      tickPositions: [
        20, 30, 40, 50, 60, 70, 80, 100, 200, 300, 400, 500, 600, 800,
        1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 10000, 20000,
      ].map((v) => Math.log10(v))
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

  constructor(
    @Inject(DOCUMENT) private document: Document,
  ) {
  }

  ngOnInit(): void {
    const ip: string = (this.document as any)['carAVServerIP'];
    fetch(`http://${ip}:8080/get`).
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
