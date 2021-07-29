import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {ImageComponent} from './image-display/image/image.component';

const appRoutes: Routes = [
    {
        path: '',
        pathMatch: 'full',
        component: ImageComponent
    }
];

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes, { relativeLinkResolution: 'legacy' })
    ]
})
export class AppRoutingModule {

}
