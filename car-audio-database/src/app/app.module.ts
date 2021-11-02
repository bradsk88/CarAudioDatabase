import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {ReactiveFormsModule} from "@angular/forms";
import {HttpClientModule} from "@angular/common/http";
import { GraphComponent } from './graph/graph.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {HighchartsChartModule} from "highcharts-angular";
import { ClaimDisplayNameComponent } from './claim-display-name/claim-display-name.component';
import { ShowFrequencyResponsesComponent } from './show-frequency-responses/show-frequency-responses.component';

@NgModule({
  declarations: [
    AppComponent,
    GraphComponent,
    ClaimDisplayNameComponent,
    ShowFrequencyResponsesComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    HighchartsChartModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
