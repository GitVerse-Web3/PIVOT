using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Zenject;

public class CameraMovement : MonoBehaviour
{
	[Inject]
	BranchManager _branchManager;

	public float movementSpeed;
	public float mouseSentity;     //鼠标灵敏度
	public Vector2 limPos;  //限制摄像机旋转的角度
	private Vector3 cameraRot;  //摄像机的三维参数
	private Vector3 cameraPos;

	// Start is called before the first frame update
	void Start()
	{
		cameraRot = new Vector3();
		cameraPos = this.transform.position;
		Input.simulateMouseWithTouches = true;
	}

	// Update is called once per frame
	private void FixedUpdate()
	{
		float mouseX = Input.GetAxis("Mouse X");
		float mouseY = Input.GetAxis("Mouse Y");

		cameraRot.x -= mouseY * mouseSentity;

		/* canmeraPosition.x赋给Roation.x，当我们向下滑动鼠标的时候，
           摄像机应该沿着X轴顺时针旋转，Roation.x的值应该变大。
           而mouseY为负数，乘以我们的灵敏度则得到了负的改变量，
           临时变量cameraPosition.x应该在原有基础上减去这个负值得到正值*/

		cameraRot.y += mouseX * mouseSentity;

		/* canmeraPosition.y赋给Roation.y，当我们向右滑动鼠标的时候，
           摄像机应该沿着y轴顺时针旋转，Roation.y的值应该变大。
           而mouseX为正数，乘以我们的灵敏度则得到了正的改变量，
           临时变量cameraPosition.x应该在原有基础上加上这个正值得到正值*/

		//cameraRot.x = Mathf.Clamp(cameraRot.x, limPos.x, limPos.y);

		/* Mathf.Clap(value，min，max)用来限制摄像机上下的旋转
           在游戏中，我们不可能360度在上下方向无限制旋转，而需要限制的值就是Roation.x的值，
           也就是cameraPosition.x的值*/

		float c = (float)_branchManager.masterHead.compressionRatio;

		if (Input.GetKey(KeyCode.X))
		{
			c = 1;
		}

		if (Input.GetKey(KeyCode.W))
			cameraPos.z += movementSpeed * c;
		if (Input.GetKey(KeyCode.S))
			cameraPos.z -= movementSpeed * c;
		if (Input.GetKey(KeyCode.A))
			cameraPos.x -= movementSpeed * c;
		if (Input.GetKey(KeyCode.D))
			cameraPos.x += movementSpeed * c;
		if (Input.GetKey(KeyCode.Space))
			cameraPos.y += movementSpeed * c;
		if (Input.GetKey(KeyCode.LeftShift))
			cameraPos.y -= movementSpeed * c;

	}

	private void LateUpdate()
	{
		//this.transform.rotation = Quaternion.Euler(0, cameraPosition.y, 0);
		//游戏人物随摄像机水平转动而转动
		transform.rotation = Quaternion.Euler(cameraRot.x, cameraRot.y, 0);
		//摄像机第一人称视角
		transform.position = cameraPos;
	}





}
