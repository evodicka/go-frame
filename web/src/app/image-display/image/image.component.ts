///<reference path="../../../../node_modules/rxjs/Observable.d.ts"/>
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ImageService } from '../service/image.service';
import { ImageInfoTo } from '../model/image-info.to';
import { Subscription } from 'rxjs/Subscription';
import { interval } from 'rxjs/observable/interval';

@Component({
  selector: 'app-image',
  templateUrl: './image.component.html',
  styleUrls: ['./image.component.css']
})
export class ImageComponent implements OnInit, OnDestroy {

  public image: ImageInfoTo;

  private subscription: Subscription;

  constructor(private imageService: ImageService) { }

  ngOnInit() {
    this.refreshImageData();
    this.schedulePeriodicRefresh();
  }

  private refreshImageData(): void {
    this.imageService.loadImageInfo().subscribe(imageData => {
      this.image = imageData;
    });
  }

  private schedulePeriodicRefresh(): void {
    this.subscription = interval(30000).subscribe(() => this.refreshImageData());
  }

  ngOnDestroy(): void {
    if (this.subscription !== null && this.subscription !== undefined) {
      this.subscription.unsubscribe();
    }
  }

}
