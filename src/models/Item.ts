import {Model, Column, Table, DataType, PrimaryKey, Unique, AutoIncrement, Length, AllowNull, Default, IsUrl, ForeignKey, BelongsTo, Scopes, DefaultScope} from "sequelize-typescript";
import {User} from './User';

type Nullable<T> = T | null;

@DefaultScope(() => ({
  attributes: { exclude: ['reserverid']}
}))

@Scopes(() => ({
  user: {
    include: [{
      model: User
    }]
  },
  reserver: {},
}))

@Table({
  modelName: "items1",
  freezeTableName: true,
  timestamps: true,
  deletedAt: false
})
export class Item extends Model<Item> {
    @PrimaryKey
    @Unique
    @AutoIncrement
    @Length({max: 11})
    @Column({
      type: DataType.INTEGER,
      field: "id"
    })
    id!: number;

    @Length({max: 255})
    @AllowNull(false)
    @Column({
      type: DataType.STRING,
      field: "name"
    })
    name!: string;

    @AllowNull(false)
    @Length({max: 255})
    @IsUrl
    @Column({
      type: DataType.STRING,
      field: "url"
    })
    url!: string;

    @AllowNull(false)
    @Default(false)
    @Column({
      type: DataType.BOOLEAN,
      field: "reserved"
    })
    reserved!: boolean;

    @Length({max: 11})
    @Default(null)
    @ForeignKey(() => User)
    @Column({
      type: DataType.INTEGER,
      field: "reserverid"
    })
    reserverid!: Nullable<number>;

    @BelongsTo(() => User)
    user?: User;

    @Default(null)
    @Column({
      type: DataType.INTEGER,
      field: "rank"
    })
    rank!: number;
}
