import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { HttpModule } from '@nestjs/axios';
import { AuthModule } from './auth/auth.module';

@Module({
  imports: [ConfigModule.forRoot(), HttpModule, AuthModule],
})
export class AppModule {}
