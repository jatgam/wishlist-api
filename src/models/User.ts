import {Model, Column, Table, DataType, PrimaryKey, Unique, AutoIncrement, Length, AllowNull, Default, IsEmail, HasMany, Scopes, DefaultScope, IsUUID, IsDate} from "sequelize-typescript";
import {Item} from './Item';

type Nullable<T> = T | null;

@DefaultScope(() => ({
  attributes: { exclude: ['hash', 'passwordreset', 'passwordResetToken', 'passwordResetExpires']}
}))

@Scopes(() => ({
  reserveditems: {
    attributes: ['id'],
    include: [{
      model: Item
    }]
  },
  auth: {

  },
  passReset: {
    attributes: {exclude: ['hash']}
  }
}))

@Table({
  modelName: "users1",
  freezeTableName: true,
  timestamps: true,
  deletedAt: false
})
export class User extends Model<User> {
    @PrimaryKey
    @Unique
    @AutoIncrement
    @Length({max: 11})
    @Column({
      type: DataType.INTEGER,
      field: "id"
    })
    id!: number;

    @Unique
    @Length({max: 30})
    @AllowNull(false)
    @Column({
      type: DataType.STRING,
      field: "username"
    })
    username!: string;

    @AllowNull(false)
    @Length({max: 255})
    @Column({
      type: DataType.STRING,
      field: "hash"
    })
    hash!: string;

    @AllowNull(false)
    @Default(false)
    @Column({
      type: DataType.BOOLEAN,
      field: "passwordreset"
    })
    passwordreset!: boolean;

    @AllowNull(true)
    @Default(null)
    @Column({
      type: DataType.STRING,
      field: "passwordResetToken"
    })
    passwordResetToken!: Nullable<string>;

    @AllowNull(true)
    @Default(null)
    @IsDate
    @Column({
      type: DataType.DATE,
      field: "passwordResetExpires"
    })
    passwordResetExpires!: Nullable<Date>;

    @AllowNull(false)
    @Length({max: 1})
    @Column({
      type: DataType.TINYINT.UNSIGNED,
      field: "userlevel"
    })
    userlevel!: number;

    @AllowNull(false)
    @Length({max: 50})
    @IsEmail
    @Column({
      type: DataType.STRING,
      field: "email"
    })
    email!: string;

    @Length({max: 30})
    @AllowNull(false)
    @Column({
      type: DataType.STRING,
      field: "firstname"
    })
    firstname!: string;

    @Length({max: 30})
    @AllowNull(false)
    @Column({
      type: DataType.STRING,
      field: "lastname"
    })
    lastname!: string;

    @HasMany(() => Item)
    reserveditems?: Item[];
}
