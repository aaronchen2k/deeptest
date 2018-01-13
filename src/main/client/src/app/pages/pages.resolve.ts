import { Injectable } from '@angular/core';
import { Router, ActivatedRoute, Resolve, ActivatedRouteSnapshot } from '@angular/router';
import {Location} from '@angular/common';
import 'rxjs/add/operator/toPromise';

import { Utils } from '../utils/utils';
import { AccountService } from '../service/account';

@Injectable()
export class PagesResolve implements Resolve<any> {
  constructor(private location: Location, private _route: ActivatedRoute, private accountService: AccountService, private router: Router) { }

  resolve(route: ActivatedRouteSnapshot) {
    let context = Utils.getOrgAndPrjId(this.location.path());

    return this.accountService.loadProfileRemote(context).toPromise().then(result => {
      return result;
    });
  }
}
