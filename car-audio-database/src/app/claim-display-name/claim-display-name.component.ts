import { Component, OnInit } from '@angular/core';
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-claim-display-name',
  templateUrl: './claim-display-name.component.html',
  styleUrls: ['./claim-display-name.component.scss']
})
export class ClaimDisplayNameComponent implements OnInit {


  constructor(private http: HttpClient) {
  }

  ngOnInit(): void {
    this.http.post(`/user/display-name/claim?displayName=${"test"}`, null).subscribe();
  }

}
