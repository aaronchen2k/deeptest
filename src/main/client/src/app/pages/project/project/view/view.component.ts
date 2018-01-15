import {Component, ViewEncapsulation, ViewChild, QueryList, Query} from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { NgModule, Pipe, OnInit, AfterViewInit }      from '@angular/core';

import {GlobalState} from '../../../../global.state';

import { CONSTANT } from '../../../../utils/constant';
import { Utils } from '../../../../utils/utils';
import {ValidatorUtils} from '../../../../validator/validator.utils';
import { RouteService } from '../../../../service/route';
import { AccountService } from '../../../../service/account';
import { ProjectService } from '../../../../service/project';

declare var jQuery;

@Component({
  selector: 'project-view',
  encapsulation: ViewEncapsulation.None,
  styleUrls: ['./view.scss'],
  templateUrl: './view.html'
})
export class ProjectView implements OnInit, AfterViewInit {
  orgId: number;
  type: string;
  id: number;
  model: any = {};

  constructor(private _state:GlobalState, private _routeService: RouteService, private _route: ActivatedRoute,
              private _projectService: ProjectService, private accountService: AccountService) {

  }
  ngOnInit() {
    let that = this;

    this.orgId = CONSTANT.CURR_ORG_ID;

    this._route.params.subscribe(params => {
      that.id = +params['id'];
      that.loadData();
    });
  }
  ngAfterViewInit() {

  }

  loadData() {
    if (CONSTANT.CURR_PRJ_ID != this.id) {
      CONSTANT.CURR_PRJ_ID = this.id;

      this._projectService.view(this.id).subscribe((json:any) => {

      });
    }
  }

}

