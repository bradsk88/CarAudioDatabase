import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {ClaimDisplayNameComponent} from "./claim-display-name/claim-display-name.component";
import {ShowFrequencyResponsesComponent} from "./show-frequency-responses/show-frequency-responses.component";

const routes: Routes = [
  {path: "app", component: ShowFrequencyResponsesComponent},
  {path: "app/display-name/claim", component: ClaimDisplayNameComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
