import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import 'reflect-metadata';

@Component({
  selector: 'app-entities',
  imports: [],
  templateUrl: './entities.component.html',
  styleUrl: './entities.component.css'
})
export class EntitiesComponent implements OnInit {
  constructor(
    private readonly httpClient: HttpClient,
  ) { }

  ngOnInit(): void {
    this.httpClient.get<any>('/entities', {}).subscribe(
      (resp) => {
        console.log(resp);
      },
      (err) => {
        console.log(err, JSON.stringify(this), Reflect.getMetadata('', this), Reflect.getMetadata('', EntitiesComponent));
      },
      () => {},
    );
    return;
  }
}
