import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

class Entity {
  Id: string = '';
  Value: string = '';
}

class MyResponse<T> {
  Results: T[] = [];
  Count: number = 0;
}

/*type Entity = {
  Id: string;
  Value: string;
}*/

@Component({
  selector: 'app-entities',
  imports: [],
  templateUrl: './entities.component.html',
  styleUrl: './entities.component.css'
})
export class EntitiesComponent implements OnInit {
  public entities: Entity[] = [];
  public count: number = 0;

  constructor(
    private readonly httpClient: HttpClient,
  ) { }

  ngOnInit(): void {
    this.httpClient.get<MyResponse<Entity>>('/entities', {responseType: 'json'}).subscribe(
      (resp) => {
        this.entities = resp.Results;
        this.count = resp.Count;
      },
      (err) => {
        console.log(err, JSON.stringify(this));
      },
      () => {},
    );
    return;
  }
}
