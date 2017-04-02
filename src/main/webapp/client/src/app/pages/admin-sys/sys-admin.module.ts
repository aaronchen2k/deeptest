import { NgModule }      from '@angular/core';
import { CommonModule }  from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { NgaModule } from '../../theme/nga.module';

import { routing }       from './sys-admin.routing';
import { CollapseModule } from 'ng2-bootstrap';
import { ModalModule } from 'ng2-bootstrap';
import { ButtonsModule } from 'ng2-bootstrap';
import { FileUploadModule } from 'ng2-file-upload';

import { SlidebarModule } from '../../components/slidebar';

import { SysAdmin } from './sys-admin.component';

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    NgaModule,
    routing,

    ModalModule,
    ButtonsModule,
    CollapseModule,
    FileUploadModule,

    SlidebarModule
  ],
  declarations: [
    SysAdmin
  ],
  providers: [

  ]
})
export default class AdminModule {}
