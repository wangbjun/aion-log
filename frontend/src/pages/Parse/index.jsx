import {Card, message, Upload} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import {InboxOutlined} from '@ant-design/icons';
import moment from "moment";

const {Dragger} = Upload;

@connect(
  state => ({
    ...state.global,
  })
)
class Parse extends React.Component {
  constructor(props) {
    super(props);
  }

  handleChange = info => {
    let resp = info.file && info.file.response
    if (resp) {
      if (resp.code !== 200) {
        message.error("失败，请重试：" + resp.msg).then();
      } else {
        message.success("成功，开始解析日志！").then();
        window.location.reload()
      }
    }
  };

  beforeUpload = async file => {
    const {dispatch} = this.props
    await dispatch({
      type: 'global/getTask',
    });
    const {taskStatus} = this.props
    if (taskStatus.isRunning) {
      message.error("当前有任务正在运行，请稍后再试！").then();
      return false
    }
    return new Promise(resolve => {
      const reader = new FileReader();
      reader.onload = function () {
        let firstTime = this.result.substring(0, 19)
        if (!firstTime) {
          message.error("此文件内容为空").then();
          return false
        }
        firstTime = firstTime.replace(new RegExp("\\.", "gm"), "-")
        if (!moment(firstTime, "YYYY-MM-DD HH:mm:ss", true).isValid()) {
          message.error("此文件格式错误").then();
          return false
        }
        if (moment(firstTime).isBefore(moment(taskStatus.lastTime))) {
          message.error("此日志已经解析过！").then();
          return false
        }
        resolve(true)
      }
      reader.readAsText(file)
    })
  }

  render() {
    const props = {
      action: '/api/v1/addTask',
      onChange: this.handleChange,
      beforeUpload: this.beforeUpload,
      multiple: false,
      maxCount: 1,
    }
    return (
      <PageContainer>
        <Card>
          <Dragger {...props}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined/>
            </p>
            <p className="ant-upload-text">点击或拖放文件上传</p>
          </Dragger>
        </Card>
      </PageContainer>
    );
  }
}

export default Parse;
