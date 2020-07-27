import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ImageInfoTo} from '../model/image-info.to';

@Injectable({
  providedIn: 'root'
})
export class ImageService {

  constructor(private httpClient: HttpClient) { }

  public loadImageInfo(): Observable<ImageInfoTo> {
    return this.httpClient.get<ImageInfoTo>('/api/image/current');
  }

}
