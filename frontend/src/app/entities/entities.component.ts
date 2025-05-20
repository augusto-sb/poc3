import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

class Entity {
  Id: string = '';
  Value: string = '';
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
  /*public readonly keys: string[] = [];
  public entities: Entity[] = [];*/
  //public readonly keys: (keyof Entity)[] = [];
  public entities: Entity[] = [];

  constructor(
    private readonly httpClient: HttpClient,
  ) {
    //console.log(Entity, Object.getOwnPropertyNames(new Entity()));
    //this.keys = Object.getOwnPropertyNames(new Entity()) as (keyof Entity)[];
  }

  ngOnInit(): void {
    this.httpClient.get<Entity[]>('/entities', {responseType: 'json'}).subscribe(
      (resp) => {
        //console.log(resp, typeof resp);
        this.entities = resp;
      },
      (err) => {
        console.log(err, JSON.stringify(this));
      },
      () => {},
    );
    return;
  }
}
