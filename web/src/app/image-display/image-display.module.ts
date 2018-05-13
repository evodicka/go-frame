import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ImageComponent } from './image/image.component';
import {ImageService} from "./service/image.service";

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [ImageComponent],
  exports: [
    ImageComponent
  ],
  providers: [
    ImageService
  ]
})
export class ImageDisplayModule { }
