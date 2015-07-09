
void Object::connectSignals()
{
    connect(this,          &Object::sendStates,
            this->otherObject,    &OtherObject::setState);

    connect(this,                     &Object::sendState,
            this->otherObject,  &OtherObject::setState);

    connect(this,                 &Object::sendState,
            this->otherObject,    &OtherObject::setState);

    connect(this,               &Object::sendState,
            this->otherObject,  &OtherObject::setState);
}
