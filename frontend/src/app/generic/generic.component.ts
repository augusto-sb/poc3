import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import 'reflect-metadata';

class Class {
  _path: string = '';
  constructor() {
    throw new Error('use a class!');
  }
}

@Component({
  selector: 'app-generic',
  imports: [],
  templateUrl: './generic.component.html',
  styleUrl: './generic.component.css'
})
export class GenericComponent<TData extends Class = Class> implements OnInit {
  public entities: TData[] = [];
  public keys: (keyof TData)[] = [];
  public count: number = 0;

  constructor(
    private readonly httpClient: HttpClient,
  ) {
    //console.log(arguments, Reflect.getMetadata, new TData())
    console.log(this)
  }

  ngOnInit(): void {
    this.httpClient.get<{results: TData[]; count: number;}>('/entities'/*AAAAA*/, {responseType: 'json'}).subscribe(
      (resp) => {
        this.entities = resp as any;
        if(this.entities.length){
          this.keys = Object.keys(this.entities[0] as any) as any;
        }
        this.count = this.entities.length;
        //this.entities = resp.results;
        //this.count = resp.count;
      },
      (err) => {
        console.log(err);alert();
      },
      () => {},
    );
    return;
  }
}