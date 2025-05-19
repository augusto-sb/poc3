import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class SessionService {
  private loggedIn: boolean = false;

  constructor(
    private readonly httpClient: HttpClient,
  ) { }

  public init(): Promise<void>{
    console.log('a')
    return new Promise((resolve, reject) => {
      this.httpClient.get('/asd', {responseType: 'text'}).subscribe(
        x => {
          console.log(x);resolve();console.log('b');this.loggedIn=x.startsWith('cookie found in session and is');
        },
        //reject,
        err => {console.log(err);console.log('c');},
        //() => {},
        resolve,
      );
    });
    console.log('d')
/*    const a = await this.httpClient.get('/asd').subscribe(
      x => {
        console.log(x);
      },
      err => {};
    );
    console.log(a)*/
  }

  public isLoggedIn(): boolean {
    return this.loggedIn;
  }
  /*public isLoggedIn(): boolean {
    return false;//(Math.random() > 0.5);
  }*/

  public login(): void {
    return;
  }

  public logout(): void {
    sessionStorage.clear();
    localStorage.clear();
    document.cookie.split('; ');
    return;
  }
}
