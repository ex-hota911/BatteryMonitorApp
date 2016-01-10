/*
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
/*
 * This code was generated by https://github.com/google/apis-client-generator/
 * (build: 2015-11-16 19:10:01 UTC)
 * on 2016-01-09 at 16:33:22 UTC 
 * Modify at your own risk.
 */

package com.appspot.icumn7abiu.batteryservice.model;

/**
 * Model definition for Device.
 *
 * <p> This is the Java data model class that specifies how to parse/serialize into the JSON that is
 * transmitted over HTTP when working with the batteryservice. For a detailed explanation see:
 * <a href="https://developers.google.com/api-client-library/java/google-http-java-client/json">https://developers.google.com/api-client-library/java/google-http-java-client/json</a>
 * </p>
 *
 * @author Google, Inc.
 */
@SuppressWarnings("javadoc")
public final class Device extends com.google.api.client.json.GenericJson {

  /**
   * The value may be {@code null}.
   */
  @com.google.api.client.util.Key("AlertThreshold")
  private java.lang.Integer alertThreshold;

  /**
   * The value may be {@code null}.
   */
  @com.google.api.client.util.Key("Batteries")
  private java.util.List<Battery> batteries;

  static {
    // hack to force ProGuard to consider Battery used, since otherwise it would be stripped out
    // see https://github.com/google/google-api-java-client/issues/543
    com.google.api.client.util.Data.nullOf(Battery.class);
  }

  /**
   * The value may be {@code null}.
   */
  @com.google.api.client.util.Key("DeviceId")
  private java.lang.String deviceId;

  /**
   * The value may be {@code null}.
   */
  @com.google.api.client.util.Key("DeviceName")
  private java.lang.String deviceName;

  /**
   * The value may be {@code null}.
   */
  @com.google.api.client.util.Key("UserId")
  private java.lang.String userId;

  /**
   * @return value or {@code null} for none
   */
  public java.lang.Integer getAlertThreshold() {
    return alertThreshold;
  }

  /**
   * @param alertThreshold alertThreshold or {@code null} for none
   */
  public Device setAlertThreshold(java.lang.Integer alertThreshold) {
    this.alertThreshold = alertThreshold;
    return this;
  }

  /**
   * @return value or {@code null} for none
   */
  public java.util.List<Battery> getBatteries() {
    return batteries;
  }

  /**
   * @param batteries batteries or {@code null} for none
   */
  public Device setBatteries(java.util.List<Battery> batteries) {
    this.batteries = batteries;
    return this;
  }

  /**
   * @return value or {@code null} for none
   */
  public java.lang.String getDeviceId() {
    return deviceId;
  }

  /**
   * @param deviceId deviceId or {@code null} for none
   */
  public Device setDeviceId(java.lang.String deviceId) {
    this.deviceId = deviceId;
    return this;
  }

  /**
   * @return value or {@code null} for none
   */
  public java.lang.String getDeviceName() {
    return deviceName;
  }

  /**
   * @param deviceName deviceName or {@code null} for none
   */
  public Device setDeviceName(java.lang.String deviceName) {
    this.deviceName = deviceName;
    return this;
  }

  /**
   * @return value or {@code null} for none
   */
  public java.lang.String getUserId() {
    return userId;
  }

  /**
   * @param userId userId or {@code null} for none
   */
  public Device setUserId(java.lang.String userId) {
    this.userId = userId;
    return this;
  }

  @Override
  public Device set(String fieldName, Object value) {
    return (Device) super.set(fieldName, value);
  }

  @Override
  public Device clone() {
    return (Device) super.clone();
  }

}