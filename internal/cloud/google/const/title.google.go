package subject_const

const ContentTitle2 = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Your Title Here</title>
    <style>
      /* Reset CSS */
      body,
      html {
        margin: 0;
        padding: 0;
        font-family: Arial, sans-serif;
        line-height: 1.6;
      }
      /* Container */
      .container {
        width: 100%;
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
        background-color: #f7f7f7;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
      }
      /* Header */
      .header {
        text-align: center;
        margin-bottom: 20px;
      }
      .header h1 {
        color: #333;
        margin: 0;
        font-size: 2.5em; /* Font size for large screens */
      }
      /* Content */
      .content {
        text-align: left;
        color: #666;
      }
      .content p {
        margin-bottom: 15px;
        font-size: 1.1em; /* Font size for large screens */
      }
      /* Button */
      .button {
        display: inline-block;
        padding: 10px 20px;
        background-color: #007bff;
        color: #fff;
        text-decoration: none;
        border-radius: 5px;
        font-size: 1em; /* Font size for large screens */
      }
      /* Responsive */
      @media screen and (max-width: 600px) {
        .container {
          border-radius: 0;
        }
        .content {
          padding: 0 10px;
        }
        .header h1 {
          font-size: 2em; /* Font size for small screens */
        }
        .content p {
          font-size: 1em; /* Font size for small screens */
        }
        .button {
          font-size: 0.9em; /* Font size for small screens */
        }
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div>
        <img
          src="https://res.cloudinary.com/df4zm1xjy/image/upload/v1713191371/feit-static/Artboard.png.png"
          alt="FEIT"
          style="width: 110px"
        />
      </div>
      <div class="header"></div>
      <div class="content">
        <p>HELLO WORLD,</p>
        <p>
          Tôi rất vui được giới thiệu với bạn về cơ hội học tập tại FEIT - một
          môi trường tuyệt vời để cải thiện kỹ năng tiếng Anh chuyên ngành và
          khám phá thế giới Công nghệ Thông tin một cách thú vị và sáng tạo.
        </p>
        <p>
          Tại FEIT, chúng tôi không chỉ xem việc học tiếng Anh và Công nghệ
          thông tin như những nhiệm vụ khô khan, mà chúng tôi còn coi đó như
          những trải nghiệm hấp dẫn, giống như việc "hack" vào một máy chủ để
          khám phá những bí ẩn mới.
        </p>
        <p>
          Chúng tôi cam kết hướng dẫn bạn cách "decode" những khía cạnh thú vị
          của cả hai lĩnh vực này một cách vui vẻ và đầy đam mê.Tại đây, bạn sẽ
          được tham gia vào các hoạt động thực hành, dự án thú vị và các buổi
          thảo luận sôi nổi, tất cả đều được thực hiện bằng tiếng Anh để giúp
          bạn nâng cao kỹ năng giao tiếp và hiểu biết chuyên môn.
        </p>
        <p>
          Nếu bạn quan tâm và muốn biết thêm thông tin, vui lòng ghé qua trang
          web của chúng tôi. Chúng tôi luôn sẵn lòng chào đón bạn!
        </p>
        <p>Trân trọng,</p>
        <p>
          Ngô Hoài Phong <br />
          Admin FEIT<br />
          gmail: hoaiphong01012002@gmail.com <br />
        </p>
        <p style="text-align: center">
          <a href="#" class="button">Learn More</a>
        </p>
      </div>
    </div>
  </body>
</html>
`
