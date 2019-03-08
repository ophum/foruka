<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class ContainerModel extends Model
{
    protected $fillable = [
        'user_id',
        'name',
        'image_id'
    ];

}
