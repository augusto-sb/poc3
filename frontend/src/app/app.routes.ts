import { Routes } from '@angular/router';

import { AboutComponent } from './about/about.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { EntitiesComponent } from './entities/entities.component';

export const routes: Routes = [
	{
		path: '',
		component: HomeComponent,
	},
	{
		path: 'about',
		component: AboutComponent,
	},
	{
		path: 'entities',
		component: EntitiesComponent,
	},
	{
		path: 'login',
		component: LoginComponent,
	},
];
