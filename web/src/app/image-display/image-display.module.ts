import { ImageComponent } from './image/image.component';
import {ImageService} from "./service/image.service";
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";

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
