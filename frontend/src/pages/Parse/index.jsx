import {Button, Card, message, Upload} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import {UploadOutlined} from '@ant-design/icons';

@connect(
  state => ({
    ...state.global,
  })
)
class Parse extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const {dispatch} = this.props
    dispatch({
      type: 'global/getTask',
    });
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

  render() {
    const {isRuning} = this.props
    const props = {
      action: '/api/v1/addTask',
      onChange: this.handleChange,
      multiple: false,
      maxCount: 1,
      showUploadList: false
    }
    return (
      <PageContainer>
        <Card>
          {isRuning ?
            <div style={{textAlign: "center", fontSize: 15}}>当前有任务正在运行，请稍后再试！</div> :
            <div style={{marginTop: "50px", textAlign: "center"}}>
              <Upload {...props}>
                <Button icon={<UploadOutlined/>} style={{fontSize: 15}}>选择文件</Button>
              </Upload>
            </div>
          }
        </Card>
      </PageContainer>
    );
  }
}

export default Parse;
